import { TicketService } from '$lib/services/ticket-service';
import { cloneDeep } from 'lodash';
import { from } from 'rxjs';
import { writable, get } from 'svelte/store';

export type BoardState = {
	initializing: boolean;
	selected?: TicketService.Board;
	boards: TicketService.Board[];
};

function defaultState(): BoardState {
	return {
		initializing: true,
		selected: undefined,
		boards: []
	};
}

const boardStore = writable<BoardState>(defaultState());

function addTicket({ board_id, ticket }: { board_id: number; ticket: TicketService.Ticket }) {
	boardStore.update((store) => {
		const boards = cloneDeep(store.boards);
		const board = boards.find((b) => b.id === board_id);
		if (board) {
			const status = board.statuses.find((s) => s.id === ticket.status_id);
			if (status) {
				status.tickets.push(ticket);
			}
		}

		store.boards = boards;

		return store;
	});
}

export const BoardStore = { ...boardStore, defaultState, addTicket };