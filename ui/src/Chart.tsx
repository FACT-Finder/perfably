import React from 'react';
import {LineChart, CartesianGrid, XAxis, YAxis, Line, Tooltip, Legend, ResponsiveContainer} from 'recharts';
import colorHash from 'color-hash';
const cHash = new colorHash();

export const Chart = (metrics: MetricRequest) => {
    const data = useMetrics(metrics);
    return (
        <ResponsiveContainer height={400}>
            <LineChart data={data}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="__key" />
                <YAxis />
                <Tooltip />
                <Legend />
                {metrics.keys.map((name) => {
                    return (
                        <Line key={name} type="monotone" dataKey={name} stroke={cHash.hex(name)} activeDot={{r: 8}} />
                    );
                })}
            </LineChart>
        </ResponsiveContainer>
    );
};

type DataPoint = Record<string, number>;
interface Metric {
    key: string;
    values: DataPoint;
}

type FlatMetric = {
    __key: string;
} & DataPoint;

interface MetricRequest {
    sort: 'asc' | 'desc';
    keys: string[];
    limit?: number;
    start?: string;
    end?: string;
    project: string;
}

const param = (name: string, value?: string | number) =>
    value === undefined ? undefined : `${encodeURIComponent(name)}=${encodeURIComponent(value)}`;

const useMetrics = ({keys, limit, start, end, sort, project}: MetricRequest): FlatMetric[] => {
    const [data, setData] = React.useState<FlatMetric[]>([]);

    React.useEffect(() => {
        const query = [
            param('start', start),
            param('end', end),
            param('limit', limit),
            param('sort', sort),
            ...keys.map((key) => param('key', key)),
        ]
            .filter((value) => value !== undefined)
            .join('&');
        fetch(`./project/${project}/value?${query}`)
            .then((res) => res.json())
            .then((data) => setData(data.map((m: Metric) => ({__key: m.key, ...m.values}))));
    }, []);
    return data;
};
