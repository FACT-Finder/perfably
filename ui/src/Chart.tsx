import React from 'react';
import {LineChart, CartesianGrid, XAxis, YAxis, Line, Tooltip, Legend, ResponsiveContainer} from 'recharts';
import colorHash from 'color-hash';
import {Unit} from './Config';
import {bestUnit, isTimeUnit} from './unit';
import {TooltipProps} from 'recharts/types/component/Tooltip';
import {Typography, Paper, Box, Table, TableBody, TableCell, TableRow, Divider} from '@mui/material';
const cHash = new colorHash();

export const Chart = (metrics: MetricRequest) => {
    const data = useMetrics(metrics);
    return (
        <ResponsiveContainer height={400}>
            <LineChart data={data.metrics}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="key" />
                <YAxis unit={data.unit} />
                <Tooltip
                    cursor={true}
                    trigger="hover"
                    content={<TooltipBody />}
                    formatter={(value: any) => (typeof value === 'number' ? data.format(value) : value)}
                />
                <Legend />
                {metrics.keys.map((name) => {
                    return (
                        <Line
                            key={name}
                            type="monotone"
                            dataKey={(item) => item.values[name]}
                            name={name}
                            stroke={cHash.hex(name)}
                            activeDot={{r: 8}}
                        />
                    );
                })}
            </LineChart>
        </ResponsiveContainer>
    );
};

type DataPoint = Record<string, number>;
type MetricPoint = Record<string, MetricValue>;
interface MetricValue {
    value: string;
    url?: string;
}
interface Metric {
    key: string;
    values: DataPoint;
    meta: MetricPoint;
}

interface ChartData {
    unit: string;
    format: (value: number) => string;
    metrics: FlatMetric[];
}

type FlatMetric = {
    values: Record<string, number>;
    meta: Record<string, {value: string}>;
    key: string;
};

interface MetricRequest {
    sort: 'asc' | 'desc';
    keys: string[];
    unit?: Unit;
    limit?: number;
    start?: string;
    end?: string;
    project: string;
}

const param = (name: string, value?: string | number) =>
    value === undefined ? undefined : `${encodeURIComponent(name)}=${encodeURIComponent(value)}`;

const useMetrics = ({keys, unit, limit, start, end, sort, project}: MetricRequest): ChartData => {
    const [data, setData] = React.useState<ChartData>({unit: '', format: (x) => x.toString(), metrics: []});

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
            .then((data) => transform(data, unit))
            .then((data) => setData(data));
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [keys, limit, start, end, sort, project]);
    return data;
};

const rescaleValues = (points: Record<string, number>, ratio: number): Record<string, number> => {
    return Object.fromEntries(Object.entries(points).map(([k, v]) => [k, ratio * v]));
};

const transform = (data: Metric[], unit?: Unit): ChartData => {
    if (isTimeUnit(unit)) {
        const max = data.reduce(
            (m, metric) => Object.values(metric.values).reduce((x, y) => Math.max(x, y), m),
            -Infinity
        );
        const {unit: targetUnit, ratio, format} = bestUnit(unit, max);
        return {
            unit: targetUnit,
            format,
            metrics: data.map(
                (m: Metric): FlatMetric => ({key: m.key, values: rescaleValues(m.values, ratio), meta: m.meta})
            ),
        };
    }
    return {
        unit: unit || '',
        format: (x) => x.toFixed(3),
        metrics: data.map((m: Metric): FlatMetric => ({key: m.key, values: m.values, meta: m.meta})),
    };
};

interface RowData {
    key: string;
    name: string;
    value: string;
    color?: string;
}
const TooltipBody = ({active, payload, label, formatter}: TooltipProps<any, any>) => {
    if (!active || !payload) {
        return null;
    }

    const meta: MetricPoint = payload?.[0]?.payload?.meta ?? {};

    const valueRows: RowData[] = payload.map((entry) => ({
        key: 'value' + entry.name,
        name: entry.name,
        color: entry.color,
        value: formatter?.(entry.value) ?? entry.value,
    }));

    const metaRows: RowData[] = Object.keys(meta)
        .sort()
        .map((key) => ({
            key: 'meta' + key,
            name: key,
            value: meta[key].value,
        }));

    const renderKey = ({name, key, color}: RowData) => (
        <Box key={key} paddingRight={1}>
            <Typography style={{color: color}}>{name}:</Typography>
        </Box>
    );
    const renderValue = ({value, key, color}: RowData) => (
        <Typography key={key} style={{color: color}} align="right">
            {value}
        </Typography>
    );

    return (
        <Paper>
            <Box padding={1}>
                <Typography variant="h5">{label}</Typography>
                <Box display="flex">
                    <Box flex={1}>{valueRows.map(renderKey)}</Box>
                    <Box>{valueRows.map(renderValue)}</Box>
                </Box>
                <Divider />
                <Box display="flex">
                    <Box flex={1}>{metaRows.map(renderKey)}</Box>
                    <Box>{metaRows.map(renderValue)}</Box>
                </Box>
            </Box>
        </Paper>
    );
};
