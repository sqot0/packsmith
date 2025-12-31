<script lang="ts">
	import { Menubar as MenubarPrimitive } from "bits-ui";
	import { cn } from "$lib/utils.js";

	let {
		ref = $bindable(null),
		class: className,
		...restProps
	}: MenubarPrimitive.TriggerProps = $props();
</script>

<MenubarPrimitive.Trigger
	bind:ref
	data-slot="menubar-trigger"
	class={cn(
		"hover:bg-accent hover:text-accent-foreground disabled:text-muted-foreground disabled:pointer-events-none flex items-center rounded-sm px-2 py-1 text-sm font-medium outline-hidden",
		className
	)}
	{...restProps}
/>
