#!/usr/bin/env node

function mostlyCapitals(line) {
	if (line
			.replace(/[^A-Za-z]/,'')
			.split(" ")
			.reduce((acc, word, _, arr) => {
				if (word[0].toLowerCase() !== word[0]) {
					acc.uppers++;
				}
				acc.total = arr.length;
				acc.isSong = ((acc.uppers/acc.total) > 0.5);
				return acc;
			}, {uppers:0}).isSong) {
		return true;
	}
	return false;
}

const readline = require('readline');

const rl = readline.createInterface({
	input: process.stdin,
	crlfDelay: Infinity
});

rl.on('line', (line) => {
	if (mostlyCapitals(line)) {
		console.log(line);
	}
});

