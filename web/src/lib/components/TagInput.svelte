<script lang="ts">
	interface Props {
		value: string[];
		suggestions?: string[];
		placeholder?: string;
		label?: string;
		onUpdate: (tags: string[]) => void;
	}

	let { value = [], suggestions = [], placeholder = 'Add tags...', label, onUpdate }: Props = $props();

	let inputValue = $state('');
	let showSuggestions = $state(false);
	let inputElement: HTMLInputElement;

	// Filter suggestions based on input and exclude already selected tags
	let filteredSuggestions = $derived(
		inputValue
			? suggestions.filter(
					(s) =>
						s.toLowerCase().includes(inputValue.toLowerCase()) &&
						!value.includes(s)
				)
			: suggestions.filter((s) => !value.includes(s))
	);

	function addTag(tag: string) {
		const trimmedTag = tag.trim();
		if (trimmedTag && !value.includes(trimmedTag)) {
			onUpdate([...value, trimmedTag]);
			inputValue = '';
			showSuggestions = false;
		}
	}

	function removeTag(tagToRemove: string) {
		onUpdate(value.filter((tag) => tag !== tagToRemove));
	}

	function handleKeyDown(e: KeyboardEvent) {
		if (e.key === 'Enter' && inputValue.trim()) {
			e.preventDefault();
			addTag(inputValue);
		} else if (e.key === 'Backspace' && !inputValue && value.length > 0) {
			// Remove last tag when backspace is pressed on empty input
			removeTag(value[value.length - 1]);
		} else if (e.key === 'Escape') {
			showSuggestions = false;
		}
	}

	function handleFocus() {
		showSuggestions = true;
	}

	function handleBlur() {
		// Delay to allow clicking on suggestions
		setTimeout(() => {
			showSuggestions = false;
		}, 200);
	}
</script>

<div class="relative">
	{#if label}
		<label class="block text-sm font-medium text-slate-300 mb-1">
			{label}
		</label>
	{/if}

	<div
		class="border border-slate-700 bg-slate-800 px-2 py-2 text-sm text-slate-100 focus-within:border-emerald-500 transition-colors"
		onclick={() => inputElement?.focus()}
	>
		<div class="flex flex-wrap gap-1.5">
			<!-- Selected Tags -->
			{#each value as tag}
				<span
					class="inline-flex items-center gap-1 bg-emerald-900/40 text-emerald-300 px-2 py-0.5 text-xs border border-emerald-700/50"
				>
					{tag}
					<button
						type="button"
						onclick={(e) => { e.stopPropagation(); removeTag(tag); }}
						class="text-emerald-400 hover:text-emerald-200 transition-colors"
						aria-label="Remove {tag}"
					>
						<svg class="h-3 w-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
						</svg>
					</button>
				</span>
			{/each}

			<!-- Input -->
			<input
				bind:this={inputElement}
				bind:value={inputValue}
				onkeydown={handleKeyDown}
				onfocus={handleFocus}
				onblur={handleBlur}
				type="text"
				{placeholder}
				class="flex-1 min-w-[120px] bg-transparent outline-none text-slate-100 placeholder-slate-500"
			/>
		</div>
	</div>

	<!-- Suggestions Dropdown -->
	{#if showSuggestions && filteredSuggestions.length > 0}
		<div
			class="absolute z-10 w-full mt-1 border border-slate-700 bg-slate-800 shadow-lg max-h-48 overflow-y-auto"
		>
			{#each filteredSuggestions as suggestion}
				<button
					type="button"
					onclick={() => addTag(suggestion)}
					class="w-full text-left px-3 py-2 text-sm text-slate-300 hover:bg-slate-700 transition-colors"
				>
					{suggestion}
				</button>
			{/each}
		</div>
	{/if}

	<!-- Create New Hint -->
	{#if inputValue.trim() && filteredSuggestions.length === 0}
		<div class="mt-1 text-xs text-slate-500">
			Press <kbd class="px-1 py-0.5 bg-slate-700 border border-slate-600 rounded text-slate-300">Enter</kbd> to create "{inputValue.trim()}"
		</div>
	{/if}
</div>
