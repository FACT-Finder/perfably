import React from 'react';
import {Config, ConfigDashboard, useConfig} from './Config';
import {
    AppBar,
    Button,
    CircularProgress,
    Container,
    Paper,
    Toolbar,
    Typography,
    Menu,
    MenuItem,
} from '@mui/material';
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
    const [project, setProject] = React.useState(() => Object.keys(config.projects)[0]);
    const [dashboard, setDashboard] = React.useState(() => config.projects[project].dashboards?.[0]);

    const projects = Object.keys(config.projects);
    const dashboards = config.projects[project].dashboards ?? [];

    const [projectAnchor, setProjectAnchor] = React.useState<null | HTMLElement>(null);
    const [dashboardAnchor, setDashboardAnchor] = React.useState<null | HTMLElement>(null);

    return (
        <Container maxWidth="md">
            <AppBar position="static">
                <Toolbar>
                    <Typography variant="h6" style={{marginRight: 10}}>
                        Perfably
                    </Typography>

                    <Button
                        variant="outlined"
                        color="inherit"
                        style={{marginRight: 10}}
                        aria-haspopup="true"
                        onClick={(event) => setProjectAnchor(event.currentTarget)}>
                        {project ?? 'no project configured'}
                    </Button>
                    <Menu
                        id="project-menu"
                        anchorEl={projectAnchor}
                        keepMounted
                        open={Boolean(projectAnchor)}
                        onClose={() => setProjectAnchor(null)}>
                        {projects.map((project) => {
                            return (
                                <MenuItem
                                    onClick={() => {
                                        setProject(project);
                                        setProjectAnchor(null);
                                    }}>
                                    {project}
                                </MenuItem>
                            );
                        })}
                    </Menu>

                    <Button
                        variant="outlined"
                        color="inherit"
                        aria-haspopup="true"
                        onClick={(event) => setDashboardAnchor(event.currentTarget)}>
                        {dashboard?.name ?? 'no dashboard configured'}
                    </Button>
                    <Menu
                        id="dashboard-menu"
                        anchorEl={dashboardAnchor}
                        keepMounted
                        open={Boolean(dashboardAnchor)}
                        onClose={() => setDashboardAnchor(null)}>
                        {dashboards.map((dashboard) => {
                            return (
                                <MenuItem
                                    onClick={() => {
                                        setDashboard(dashboard);
                                        setDashboardAnchor(null);
                                    }}>
                                    {dashboard.name}
                                </MenuItem>
                            );
                        })}
                    </Menu>
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
