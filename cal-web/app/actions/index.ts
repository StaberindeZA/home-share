'use server'

import { rows } from "../constants";

function convertTimeToUTCISO(timeString: string) {
	const [hours, minutes, seconds] = timeString.split(':').map(Number);
	const date = new Date();
	date.setHours(hours, minutes, seconds, 0);
	const isoString = date.toISOString();
	return isoString.replace(/\.\d{3}/, '');
}
const wait = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

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

export async function listEntries() {
	return rows;
}
