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
    const [filterAnchor, setFilterAnchor] = React.useState<null | HTMLElement>(null);
    const ids = useIds(project);
    const [range, setRange] = React.useState<[number, number] | undefined>(undefined);
    React.useEffect(() => {
        if (range === undefined && ids.length !== 0) {
            setRange([Math.max(0, ids.length - 1 - 50), ids.length - 1]);
        }
    }, [ids, range]);

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
                            Dashboard: {dashboard?.name ?? 'not configured'}
                        </Button>
                        <Button
                            onClick={(event) => setFilterAnchor(event.currentTarget)}
                            endIcon={<KeyboardArrowDownIcon />}
                        >
                            Filter: {range === undefined ? 'unknown' : `${ids[range[0]]} to ${ids[range[1]]}`}
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
                                        setProject(project);
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
                        {dashboards.map((dashboard) => {
                            return (
                                <MenuItem
                                    key={dashboard.name}
                                    onClick={() => {
                                        setDashboard(dashboard);
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
                        {range !== undefined ? (
                            <ClickAwayListener onClickAway={() => setFilterAnchor(null)}>
                                <Box paddingX={4} paddingY={2}>
                                    <Slider
                                        valueLabelFormat={(value) => ids[value]}
                                        min={0}
                                        max={ids.length - 1}
                                        value={range}
                                        style={{paddingTop: 50}}
                                        valueLabelDisplay="on"
                                        onChange={(_, value) => setRange(value as [number, number])}
                                    />
                                </Box>
                            </ClickAwayListener>
                        ) : undefined}
                    </Popover>
                </Toolbar>
            </AppBar>
            <Box marginTop="100px" paddingX={3}>
                {dashboard && range !== undefined ? (
                    <Dashboard project={project} dashboard={dashboard} range={[ids[range[0]], ids[range[1]]]} />
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
