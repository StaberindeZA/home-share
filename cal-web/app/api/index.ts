import { API_URL } from "../constants";
import { convertTimeToUTCISO } from "../utils";

export async function createEntry(userId: number, startTime: string, endTime: string) {
	const startDateString = convertTimeToUTCISO(startTime)
	const endDateString = convertTimeToUTCISO(endTime)
}

export async function fetchListEntries(userId: number) {
	const startDateString = convertTimeToUTCISO("00:00:00")
	const endDateString = convertTimeToUTCISO("23:59:59")
	const params = {
		userId: `${userId}`,
		start: startDateString,
		end: endDateString
	}
	const url = new URL(`${API_URL}/v1/entry`);
	url.search = new URLSearchParams(params).toString();
	try {
		const results = await fetch(url);
		return results.json();
	} catch (err) {
		console.error("Something bad happend", err)
		return []
	}
}


