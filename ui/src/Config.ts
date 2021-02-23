import React from 'react';

export interface Config {
    projects: Record<string, ConfigProject>;
}
export interface ConfigProject {
    name: string;
    layers: string[];
    dashboards?: ConfigDashboard[];
}
export interface ConfigDashboard {
    name: string;
    charts?: ConfigChart[];
}
export interface ConfigChart {
    name: string;
    metrics?: string[];
}

export const useConfig = (): Config | undefined => {
    const [data, setData] = React.useState<Config | undefined>();

    React.useEffect(() => {
        fetch(`./config`)
            .then((res) => res.json())
            .then(setData);
    }, []);
    return data;
};
