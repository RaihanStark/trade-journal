<script lang="ts">
	import { authStore } from '$lib/stores/auth.svelte';
	import { z } from 'zod';

	let email = $state('');
	let password = $state('');
	let confirmPassword = $state('');
	let error = $state('');
	let isLoading = $state(false);
	let errors = $state<Record<string, string>>({});

	// Zod validation schema
	const registerSchema = z.object({
		email: z.string()
			.min(1, 'Email is required')
			.email('Please enter a valid email address'),
		password: z.string()
			.min(8, 'Password must be at least 8 characters'),
		confirmPassword: z.string()
			.min(1, 'Please confirm your password')
	}).refine((data) => data.password === data.confirmPassword, {
		message: "Passwords do not match",
		path: ["confirmPassword"],
	});

	async function handleSubmit(e: Event) {
		e.preventDefault();
		error = '';
		errors = {};

		// Validate form data
		const result = registerSchema.safeParse({
			email,
			password,
			confirmPassword
		});

		if (!result.success) {
			// Map Zod errors to field errors
			result.error.issues.forEach((err) => {
				if (err.path[0]) {
					errors[err.path[0] as string] = err.message;
				}
			});
			return;
		}

		isLoading = true;

		try {
			const err = await authStore.register(email, password);
			if (err) {
				error = err;
			}
		} catch (err) {
			error = 'An unexpected error occurred';
		} finally {
			isLoading = false;
		}
	}

	function validateField(field: string, value: any) {
		// Validate single field
		try {
			// For confirmPassword, we need to validate the whole object
			if (field === 'confirmPassword') {
				registerSchema.parse({ email, password, confirmPassword });
				const { confirmPassword: _, ...rest } = errors;
				errors = rest;
			} else {
				const fieldSchema = registerSchema.shape[field as keyof typeof registerSchema.shape];
				if (fieldSchema) {
					fieldSchema.parse(value);
					// Clear error if validation passes
					const { [field]: _, ...rest } = errors;
					errors = rest;
				}
			}
		} catch (err) {
			if (err instanceof z.ZodError) {
				const fieldError = err.issues.find(issue => issue.path.includes(field));
				if (fieldError) {
					errors[field] = fieldError.message;
				}
			}
		}
	}
</script>

<div class="flex h-full items-center justify-center bg-slate-950">
	<div class="w-full max-w-md">
		<!-- Header -->
		<div class="mb-8 text-center">
			<div class="mb-4 flex justify-center">
				<div class="flex h-12 w-12 items-center justify-center bg-slate-800 font-mono text-sm font-bold text-emerald-400">
					FX
				</div>
			</div>
			<h1 class="mb-2 text-2xl font-bold text-slate-100">FOREX JOURNAL</h1>
			<p class="text-sm text-slate-500">Create your account</p>
		</div>

		<!-- Register Form -->
		<div class="border border-slate-800 bg-slate-900 p-8">
			<form onsubmit={handleSubmit} class="space-y-6">
				{#if error}
					<div class="border border-red-800 bg-red-900/20 px-4 py-3 text-sm text-red-400">
						{error}
					</div>
				{/if}

				<!-- Email Field -->
				<div>
					<label for="email" class="mb-2 block text-xs font-bold uppercase text-slate-400">
						Email
					</label>
					<input
						type="email"
						id="email"
						bind:value={email}
						onblur={() => validateField('email', email)}
						class="w-full border bg-slate-950 px-4 py-3 text-sm text-slate-100 transition-colors placeholder:text-slate-600 focus:outline-none {errors.email
							? 'border-red-500 bg-red-900/10 focus:border-red-500'
							: 'border-slate-700 focus:border-emerald-500'}"
						placeholder="your@email.com"
					/>
					{#if errors.email}
						<p class="mt-1 text-xs text-red-400">{errors.email}</p>
					{/if}
				</div>

				<!-- Password Field -->
				<div>
					<label for="password" class="mb-2 block text-xs font-bold uppercase text-slate-400">
						Password
					</label>
					<input
						type="password"
						id="password"
						bind:value={password}
						onblur={() => validateField('password', password)}
						class="w-full border bg-slate-950 px-4 py-3 text-sm text-slate-100 transition-colors placeholder:text-slate-600 focus:outline-none {errors.password
							? 'border-red-500 bg-red-900/10 focus:border-red-500'
							: 'border-slate-700 focus:border-emerald-500'}"
						placeholder="Minimum 8 characters"
					/>
					{#if errors.password}
						<p class="mt-1 text-xs text-red-400">{errors.password}</p>
					{/if}
				</div>

				<!-- Confirm Password Field -->
				<div>
					<label for="confirm-password" class="mb-2 block text-xs font-bold uppercase text-slate-400">
						Confirm Password
					</label>
					<input
						type="password"
						id="confirm-password"
						bind:value={confirmPassword}
						onblur={() => validateField('confirmPassword', confirmPassword)}
						class="w-full border bg-slate-950 px-4 py-3 text-sm text-slate-100 transition-colors placeholder:text-slate-600 focus:outline-none {errors.confirmPassword
							? 'border-red-500 bg-red-900/10 focus:border-red-500'
							: 'border-slate-700 focus:border-emerald-500'}"
						placeholder="Re-enter your password"
					/>
					{#if errors.confirmPassword}
						<p class="mt-1 text-xs text-red-400">{errors.confirmPassword}</p>
					{/if}
				</div>

				<!-- Submit Button -->
				<button
					type="submit"
					disabled={isLoading}
					class="w-full bg-emerald-600 px-4 py-3 text-sm font-bold uppercase text-white transition-colors hover:bg-emerald-700 disabled:cursor-not-allowed disabled:opacity-50"
				>
					{isLoading ? 'Creating account...' : 'Create Account'}
				</button>
			</form>

			<!-- Login Link -->
			<div class="mt-6 border-t border-slate-800 pt-6 text-center">
				<p class="text-sm text-slate-500">
					Already have an account?
					<a href="/login" class="font-medium text-emerald-400 hover:text-emerald-300">
						Sign in
					</a>
				</p>
			</div>
		</div>
	</div>
</div>
