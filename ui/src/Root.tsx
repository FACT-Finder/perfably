import React from 'react';
import {Config, ConfigDashboard, useConfig} from './Config';
import {AppBar, Button, CircularProgress, Container, Paper, Toolbar, Typography} from '@material-ui/core';
import {Chart} from './Chart';

export const Root = () => {
    const config = useConfig();

    if (config === undefined) {
        return <CircularProgress />;
    }

    if (Object.keys(config.projects).length === 0) {
        return <>no projects configured :/</>;
    }

    return <WithConfig config={config} />;
};

const WithConfig = ({config}: {config: Config}) => {
    const [project] = React.useState(() => Object.keys(config.projects)[0]);
    const [dashboard] = React.useState(() => config.projects[project].dashboards?.[0]);

    return (
        <Container maxWidth="md">
            <AppBar position="static">
                <Toolbar>
                    <Typography variant="h6" style={{marginRight: 10}}>
                        Perfably
                    </Typography>
                    <Button variant="outlined" color="inherit" style={{marginRight: 10}}>
                        {project}
                    </Button>
                    {dashboard === undefined ? (
                        'no dashboards configured'
                    ) : (
                        <Button variant="outlined" color="inherit">
                            {dashboard.name}
                        </Button>
                    )}
                </Toolbar>
            </AppBar>
            {dashboard ? <Dashboard project={project} dashboard={dashboard} /> : undefined}
        </Container>
    );
};

const Dashboard = ({project, dashboard}: {project: string; dashboard: ConfigDashboard}) => {
    const charts = dashboard.charts ?? [];

    return (
        <>
            {charts.map((chart) => {
                return (
                    <Paper elevation={5} style={{marginTop: 10, padding: 10}}>
                        <Typography variant="h4" align="center">
                            {chart.name}
                        </Typography>
                        <Chart sort="asc" keys={chart.metrics ?? []} project={project} />
                    </Paper>
                );
            })}
        </>
    );
};
