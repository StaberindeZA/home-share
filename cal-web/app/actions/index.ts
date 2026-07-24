'use server'

import { revalidatePath } from "next/cache";
import { createEntry, deleteEntry, fetchListEntries } from "../api";
import { rows, startIndexMap, users } from "../constants";
import { convertTimeToUTCISO, convertToLocalTime, wait } from "../utils";

export async function addEntry(userId: number, startTime: string, endTime: string) {
	const startDateString = convertTimeToUTCISO(startTime);
	const endDateString = convertTimeToUTCISO(endTime);
	await createEntry(userId, startDateString, endDateString)

	revalidatePath('/')
}

export async function removeEntry(entryId: number) {
	await deleteEntry(entryId);

	revalidatePath('/')
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
