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
    ButtonGroup,
    Box,
    ClickAwayListener,
    Popper,
    Autocomplete,
    TextField,
} from '@mui/material';
import KeyboardArrowDownIcon from '@mui/icons-material/KeyboardArrowDown';
import {Chart} from './Chart';
import {useIds} from './ids';
import {useUrlChangableState} from './state';
import {useDebounce} from './useDebounce';

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
    const [autoCompleteFrom, setAutoCompleteFrom] = React.useState('');
    const [autoCompleteTo, setAutoCompleteTo] = React.useState('');

    const [{project, dashboard, filter}, setState] = useUrlChangableState();
    React.useEffect(() => {
        if (filter) {
            setAutoCompleteFrom(filter[0]);
            setAutoCompleteTo(filter[1]);
        }
    }, [filter, setAutoCompleteTo, setAutoCompleteFrom]);

    const setFilter = React.useCallback(
        (from: string, to: string) => {
            setAutoCompleteFrom(from);
            setAutoCompleteTo(to);
            setState((c) => ({...c, filter: [from, to]}));
        },
        [setState, setAutoCompleteTo, setAutoCompleteFrom]
    );

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

    const debouncedFilter = useDebounce(filter, 200);
    const projects = Object.keys(config.projects);
    const dashboards = project ? config.projects?.[project]?.dashboards ?? [] : [];

    const [projectAnchor, setProjectAnchor] = React.useState<null | HTMLElement>(null);
    const [dashboardAnchor, setDashboardAnchor] = React.useState<null | HTMLElement>(null);
    const [filterAnchor, setFilterAnchor] = React.useState<null | HTMLElement>(null);
    const ids = useIds(project);

    React.useEffect(() => {
        if (ids.length > 0 && (filter === undefined || filter.some((v) => !ids.includes(v)))) {
            setFilter(ids[Math.max(0, ids.length - 1 - 50)], ids[ids.length - 1]);
        }
    }, [ids, filter, setFilter]);

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
                    <Popper
                        placement="top-start"
                        style={{maxWidth: 500, zIndex: 1200}}
                        anchorEl={filterAnchor}
                        open={!!filterAnchor}
                    >
                        {filter !== undefined ? (
                            <ClickAwayListener onClickAway={() => setFilterAnchor(null)}>
                                <Paper>
                                    <Box padding={2} display="flex" alignItems="center">
                                        <Autocomplete
                                            size="small"
                                            options={ids}
                                            value={autoCompleteFrom}
                                            style={{width: 150}}
                                            freeSolo={false}
                                            onChange={(_, fromFilter, reason) => {
                                                if (reason === 'selectOption') {
                                                    const fromIndex = ids.indexOf(fromFilter!);
                                                    let toFilter = filter[1];
                                                    if (fromIndex > ids.indexOf(toFilter)) {
                                                        toFilter = ids[Math.min(ids.length - 1, fromIndex + 50)];
                                                    }
                                                    setFilter(fromFilter!, toFilter);
                                                } else {
                                                    setAutoCompleteFrom(fromFilter!);
                                                }
                                            }}
                                            renderInput={(params) => <TextField {...params} />}
                                        />
                                        <Box paddingX={2}>
                                            <Typography>to</Typography>
                                        </Box>
                                        <Autocomplete
                                            size="small"
                                            options={ids}
                                            freeSolo={false}
                                            value={autoCompleteTo}
                                            style={{width: 150}}
                                            onChange={(_, toFilter, reason) => {
                                                if (reason === 'selectOption') {
                                                    const toIndex = ids.indexOf(toFilter!);
                                                    let fromFilter = filter[0];
                                                    if (toIndex < ids.indexOf(fromFilter)) {
                                                        fromFilter = ids[Math.max(0, toIndex - 50)];
                                                    }
                                                    setFilter(fromFilter, toFilter!);
                                                } else {
                                                    setAutoCompleteTo(toFilter!);
                                                }
                                            }}
                                            renderInput={(params) => <TextField {...params} />}
                                        />
                                    </Box>
                                </Paper>
                            </ClickAwayListener>
                        ) : undefined}
                    </Popper>
                </Toolbar>
            </AppBar>
            <Box marginTop="100px" paddingX={3}>
                {project !== undefined &&
                dashboard !== undefined &&
                dashboards[dashboard] &&
                debouncedFilter !== undefined ? (
                    <Dashboard project={project} dashboard={dashboards[dashboard]} range={debouncedFilter} />
                ) : undefined}
            </Box>
        </Box>
    );
};

const Dashboard = React.memo(
    ({project, dashboard, range}: {project: string; dashboard: ConfigDashboard; range: [string, string]}) => {
        const charts = dashboard.charts ?? [];

        return (
            <>
                {charts.map((chart) => {
                    return (
                        <Paper key={chart.name} elevation={5} style={{marginTop: 10, padding: 10}}>
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
    }
);
