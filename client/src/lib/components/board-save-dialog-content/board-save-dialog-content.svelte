<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index';
	import { Input } from '$lib/components/ui/input/index';
	import { Label } from '$lib/components/ui/label/index';
	import { TicketService } from '$lib/services/ticket-service';
	import { Button } from '$lib/components/ui/button/index';
	import { BoardStore } from '$lib/stores/board';
	import { AlertStore } from '$lib/stores/alert';
	import { DialogStore } from '$lib/stores/dialog';
	import { onMount } from 'svelte';

	export let model = {
		id: 0,
		title: ''
	};

	onMount(() => {
		data.title = model.title;
	});

	let data = {
		title: ''
	};

	async function handleSubmit() {
		if (model.id > 0) {
			try {
				const { data: board } = await TicketService.updateBoard({ board_id: model.id }, data);

				BoardStore.updateBoard({ board });
				DialogStore.close();
			} catch (error: any) {
				AlertStore.error(error);
			}
			return;
		}

		try {
			const { data: board } = await TicketService.createBoard(data);

			BoardStore.addBoard({ board });
			DialogStore.close();
		} catch (error: any) {
			AlertStore.error(error);
		}
	}
</script>

<Dialog.Content>
	<Dialog.Header>
		<Dialog.Title>{model.id > 0 ? 'Update' : 'Create'} board</Dialog.Title>
	</Dialog.Header>
	<div class="grid gap-4 py-4">
		<div class="grid gap-4 py-4">
			<div class="grid items-center grid-cols-4 gap-4">
				<Label for="name" class="text-right">Title</Label>
				<Input id="name" bind:value={data.title} class="col-span-3" />
			</div>
		</div>
	</div>
	<Dialog.Footer>
		<Button type="submit" on:click={handleSubmit}>Create</Button>
	</Dialog.Footer>
</Dialog.Content>
