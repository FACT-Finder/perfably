import React from 'react';
import {Config, ConfigDashboard, useConfig} from './Config';
import {
    AppBar,
    Button,
    CircularProgress,
    Paper,
    Toolbar,
    Typography,
    Menu,
    MenuItem,
    Slider,
    ButtonGroup,
    Box,
    Popover,
    ClickAwayListener,
} from '@mui/material';
import KeyboardArrowDownIcon from '@mui/icons-material/KeyboardArrowDown';
import {Chart} from './Chart';
import {useIds} from './ids';
import {useUrlChangableState} from './state';

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
    const [{project, dashboard, filter}, setState] = useUrlChangableState();
    React.useEffect(() => {
        const projects = Object.keys(config.projects);
        let nextState = {project, dashboard};
        if (nextState.project === undefined || !projects.includes(nextState.project as string)) {
            if (projects.length === 0) {
                return;
            }
            nextState.project = projects[0];
        }
        const availableDashboards = config.projects[nextState.project!!].dashboards ?? [];

        if (dashboard === undefined || availableDashboards[dashboard] === undefined) {
            if (availableDashboards.length === 0) {
                return;
            }
            nextState.dashboard = 0;
        }
        setState((c) => ({...c, ...nextState}));
    }, [project, dashboard, setState, config]);

    const projects = Object.keys(config.projects);
    const dashboards = project ? config.projects?.[project]?.dashboards ?? [] : [];

    const [projectAnchor, setProjectAnchor] = React.useState<null | HTMLElement>(null);
    const [dashboardAnchor, setDashboardAnchor] = React.useState<null | HTMLElement>(null);
    const [filterAnchor, setFilterAnchor] = React.useState<null | HTMLElement>(null);
    const ids = useIds(project);

    const start = ids.findIndex((id) => id === filter?.[0]);
    const end = ids.findIndex((id) => id === filter?.[1]);
    React.useEffect(() => {
        if (ids.length > 0 && (start === -1 || end === -1)) {
            setState((c) => ({...c, filter: [ids[Math.max(0, ids.length - 1 - 50)], ids[ids.length - 1]]}));
        }
    }, [ids, start, end, setState]);
    const indexFilter: [number, number] | undefined = start !== -1 && end !== -1 ? [start, end] : undefined;

    return (
        <Box>
            <AppBar position="fixed">
                <Toolbar>
                    <Typography variant="h6" style={{marginRight: 30}}>
                        Perfably
                    </Typography>
                    <ButtonGroup variant="outlined" aria-label="outlined primary button group" color="inherit">
                        <Button
                            onClick={(event) => setProjectAnchor(event.currentTarget)}
                            endIcon={<KeyboardArrowDownIcon />}
                        >
                            Project: {project ?? 'not configured'}
                        </Button>
                        <Button
                            onClick={(event) => setDashboardAnchor(event.currentTarget)}
                            endIcon={<KeyboardArrowDownIcon />}
                        >
                            Dashboard: {dashboards[dashboard ?? -1]?.name ?? 'not configured'}
                        </Button>
                        <Button
                            onClick={(event) => setFilterAnchor(event.currentTarget)}
                            endIcon={<KeyboardArrowDownIcon />}
                        >
                            Filter: {filter === undefined ? 'unknown' : `${filter[0]} to ${filter[1]}`}
                        </Button>
                    </ButtonGroup>
                    <Menu
                        id="project-menu"
                        anchorEl={projectAnchor}
                        keepMounted
                        open={Boolean(projectAnchor)}
                        onClose={() => setProjectAnchor(null)}
                    >
                        {projects.map((project) => {
                            return (
                                <MenuItem
                                    key={project}
                                    onClick={() => {
                                        setState((c) => ({...c, project}));
                                        setProjectAnchor(null);
                                    }}
                                >
                                    {project}
                                </MenuItem>
                            );
                        })}
                    </Menu>
                    <Menu
                        id="dashboard-menu"
                        anchorEl={dashboardAnchor}
                        keepMounted
                        open={Boolean(dashboardAnchor)}
                        onClose={() => setDashboardAnchor(null)}
                    >
                        {dashboards.map((dashboard, index) => {
                            return (
                                <MenuItem
                                    key={dashboard.name}
                                    onClick={() => {
                                        setState((c) => ({...c, dashboard: index}));
                                        setDashboardAnchor(null);
                                    }}
                                >
                                    {dashboard.name}
                                </MenuItem>
                            );
                        })}
                    </Menu>
                    <Popover
                        anchorEl={filterAnchor}
                        open={!!filterAnchor}
                        anchorOrigin={{
                            vertical: 'bottom',
                            horizontal: 'center',
                        }}
                        transformOrigin={{
                            vertical: 'top',
                            horizontal: 'center',
                        }}
                        PaperProps={{
                            style: {width: '100%', maxWidth: 500},
                        }}
                    >
                        {indexFilter !== undefined ? (
                            <ClickAwayListener onClickAway={() => setFilterAnchor(null)}>
                                <Box paddingX={4} paddingY={2}>
                                    <Slider
                                        valueLabelFormat={(value) => ids[value]}
                                        min={0}
                                        max={ids.length - 1}
                                        value={indexFilter}
                                        style={{paddingTop: 50}}
                                        valueLabelDisplay="on"
                                        onChange={(_, value) => {
                                            const [left, right] = value as [number, number];
                                            setState((c) => ({...c, filter: [ids[left], ids[right]]}));
                                        }}
                                    />
                                </Box>
                            </ClickAwayListener>
                        ) : undefined}
                    </Popover>
                </Toolbar>
            </AppBar>
            <Box marginTop="100px" paddingX={3}>
                {project !== undefined && dashboard !== undefined && dashboards[dashboard] && filter !== undefined ? (
                    <Dashboard project={project} dashboard={dashboards[dashboard]} range={filter} />
                ) : undefined}
            </Box>
        </Box>
    );
};

const Dashboard = ({
    project,
    dashboard,
    range,
}: {
    project: string;
    dashboard: ConfigDashboard;
    range: [string, string];
}) => {
    const charts = dashboard.charts ?? [];

    return (
        <>
            {charts.map((chart) => {
                return (
                    <Paper elevation={5} style={{marginTop: 10, padding: 10}}>
                        <Typography variant="h4" align="center">
                            {chart.name}
                        </Typography>
                        <Chart
                            sort="asc"
                            keys={chart.metrics ?? []}
                            unit={chart.unit}
                            project={project}
                            start={range[0]}
                            end={range[1]}
                        />
                    </Paper>
                );
            })}
        </>
    );
};
