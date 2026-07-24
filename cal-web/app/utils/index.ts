export function convertTimeToUTCISO(timeString: string) {
	const [hours, minutes, seconds] = timeString.split(':').map(Number);
	const date = new Date();
	date.setHours(hours, minutes, seconds, 0);
	const isoString = date.toISOString();
	return isoString.replace(/\.\d{3}/, '');
}

export function convertToLocalTime(date: Date, timeZone: string) {
	return date.toLocaleString('en-US', {
		timeZone,
		hour: '2-digit',
		minute: '2-digit',
		second: '2-digit',
		hour12: false
	});
}

export const wait = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));


