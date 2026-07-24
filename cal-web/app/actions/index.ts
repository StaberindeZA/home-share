'use server'

import { fetchListEntries } from "../api";
import { rows, startIndexMap, users } from "../constants";
import { convertTimeToUTCISO, convertToLocalTime, wait } from "../utils";

export async function addEntry(userId: number, startTime: string, endTime: string) {
	console.log("Add:", userId, startTime, endTime)
	const startDateString = convertTimeToUTCISO(startTime);
	const endDateString = convertTimeToUTCISO(endTime);
	await wait(500);
	console.log("Add done", startDateString, endDateString)
}

export async function deleteEntry(entryId: number) {
	console.log("Delete:", entryId)
	await wait(500);
	console.log("Delete done")
}

export async function listEntries(timeZone: string) {
	const currentRows = structuredClone(rows)
	const entries = await Promise.all([
		fetchListEntries(users[0].id),
		fetchListEntries(users[1].id)
	])


	entries.forEach((entry, userIndex) => {
		entry.forEach((e: any) => {
			const startTime = convertToLocalTime(new Date(e.start), timeZone);
			const startIndex = startIndexMap.get(startTime);

			if (startIndex) {
				currentRows[startIndex].entryIds[userIndex] = e.id
			} else {
				console.error("Could not find index for:", startTime)
			}

		})
	})

	return currentRows;
}
