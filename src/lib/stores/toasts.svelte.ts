export interface Toast {
	id: string;
	type: 'success' | 'error' | 'info' | 'warning';
	message: string;
	duration?: number;
}

function createToastStore() {
	let toasts = $state<Toast[]>([]);

	function add(type: Toast['type'], message: string, duration = 4000) {
		const id = crypto.randomUUID();
		toasts.push({ id, type, message, duration });
		if (duration > 0) {
			setTimeout(() => remove(id), duration);
		}
	}

	function remove(id: string) {
		toasts = toasts.filter((t) => t.id !== id);
	}

	return {
		get items() { return toasts; },
		success: (msg: string) => add('success', msg),
		error: (msg: string) => add('error', msg, 6000),
		info: (msg: string) => add('info', msg),
		warning: (msg: string) => add('warning', msg),
		remove
	};
}

export const toastStore = createToastStore();
