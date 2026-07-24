import { API_URL } from "../constants";
import { convertTimeToUTCISO } from "../utils";

export async function createEntry(userId: number, start: string, end: string) {
	const data = {
		userId: `${userId}`,
		start,
		end,
	}
	const url = new URL(`${API_URL}/v1/entry`);
	try {
		const response = await fetch(url, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(data)
		});

		if (!response.ok) {
			throw new Error(`HTTP error! Status: ${response.status}`);
		}

		return await response.json();
	} catch (error) {
		console.error('Error:', error);
	}
}

export async function deleteEntry(entryId: number) {
	const url = new URL(`${API_URL}/v1/entry/${entryId}`);
	try {
		const response = await fetch(url, {
			method: 'DELETE',
		});

		if (!response.ok) {
			throw new Error(`HTTP error! Status: ${response.status}`);
		}

		return await response.json();
	} catch (error) {
		console.error('Error:', error);
	}
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
		const response = await fetch(url);
		if (!response.ok) {
			throw new Error(`HTTP error! Status: ${response.status}`);
		}
		return response.json();
	} catch (err) {
		console.error("Something bad happend", err)
		return []
	}
}


