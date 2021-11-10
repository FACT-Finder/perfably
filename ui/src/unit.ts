const ordered = ['ns', 'us', 'ms', 's', 'm', 'h'] as const;
export type TimeUnit = typeof ordered[number];

export const isTimeUnit = (unit?: string): unit is TimeUnit => ordered.includes(unit as TimeUnit);

interface TimeUnitMeta {
    factor: number;
    format: (x: number) => string;
}

const meta = (unit: TimeUnit): TimeUnitMeta => {
    switch (unit) {
        case 'ns':
            return {factor: 1e-9, format: (x) => x.toFixed(3) + unit};
        case 'us':
            return {factor: 1e-6, format: (x) => x.toFixed(3) + unit};
        case 'ms':
            return {factor: 1e-3, format: (x) => x.toFixed(3) + unit};
        case 's':
            return {factor: 1, format: (x) => x.toFixed(3) + unit};
        case 'm':
            return {factor: 60, format: (x) => x.toFixed(2) + unit};
        case 'h':
            return {factor: 3600, format: (x) => x.toFixed(2) + unit};
    }
};

const ratio = (from: TimeUnit, to: TimeUnit): number => {
    return meta(from).factor / meta(to).factor;
};

export const bestUnit = (from: TimeUnit, x: number): {unit: TimeUnit; ratio: number; format: (x: number) => string} => {
    const target = ordered.reduce((prev, current) => (ratio(from, current) * x < 5 ? prev : current), from);
    return {unit: target, ratio: ratio(from, target), format: meta(target).format};
};
