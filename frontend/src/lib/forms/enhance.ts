import type { SubmitFunction } from '@sveltejs/kit';

export function createFormEnhance(setLoading: (loading: boolean) => void): SubmitFunction {
	return () => {
		setLoading(true);
		return async ({ result, update }) => {
			setLoading(false);
			if (result.type === 'redirect') {
				await update();
			} else {
				await update({ reset: false });
			}
		};
	};
}
