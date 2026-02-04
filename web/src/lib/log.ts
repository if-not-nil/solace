/* eslint-disable @typescript-eslint/no-explicit-any */
const LOG_LEVELS = {
	debug: 0,
	info: 1,
	warn: 2,
	error: 3
};

const currentLogLevel = import.meta.env.DEV ? LOG_LEVELS.debug : LOG_LEVELS.warn;

export function logger(place: string) {
	const formatMessage = (level: string, msg: any) => {
		return `[${place}][${level}] - ${JSON.stringify(msg)}`;
	};

	return {
		info: (...args: any[]) => {
			if (currentLogLevel <= LOG_LEVELS.info) {
				console.info(formatMessage('info', args));
			}
		},
		err: (...args: any[]) => {
			if (currentLogLevel <= LOG_LEVELS.error) {
				console.error(formatMessage('error', args));
			}
		},
		warn: (...args: any[]) => {
			if (currentLogLevel <= LOG_LEVELS.warn) {
				console.warn(formatMessage('warn', args));
			}
		},
		debug: (...args: any[]) => {
			if (currentLogLevel <= LOG_LEVELS.debug) {
				console.debug(formatMessage('debug', args));
			}
		}
	};
}
