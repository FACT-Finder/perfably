import React, {SetStateAction} from 'react';

export const getStateFromURL = (search: string): State => {
    const pairs = search.slice(1).split('&');
    const dashboard = pairs.find((param) => param.startsWith('dashboard='))?.split('=')[1];
    const project = pairs.find((param) => param.startsWith('project='))?.split('=')[1] ?? '';
    const filter = pairs.find((param) => param.startsWith('filter='))?.split('=')[1] ?? '';
    const [left, right] = filter.split(':');
    return {
        dashboard: dashboard ? +decodeURIComponent(dashboard) : undefined,
        project: project ? decodeURIComponent(project) : undefined,
        filter: left && right ? [left, right] : undefined,
    };
};

export interface State {
    dashboard?: number;
    project?: string;
    filter?: [string, string];
}
export type NavigateNote = (note: string) => void;
export type SetState = React.Dispatch<SetStateAction<State>>;

export const useUrlChangableState = (): [State, SetState] => {
    const [state, setState] = React.useState<State>(() => getStateFromURL(window.location.search));
    React.useEffect(() => {
        const onChange = (): void => setState(getStateFromURL(window.location.search));
        window.addEventListener('popstate', onChange);
        return () => window.removeEventListener('popstate', onChange);
    }, [setState]);

    const setStateAndUrl = React.useCallback(
        (stateF: SetStateAction<State>) => {
            setState((old) => {
                const newState = typeof stateF === 'function' ? stateF(old) : stateF;

                if (
                    old.dashboard === newState.dashboard &&
                    old.project === newState.project &&
                    old.filter?.[0] === newState.filter?.[0] &&
                    old.filter?.[1] === newState.filter?.[1]
                ) {
                    return old;
                }
                const params: string[] = [];
                if (newState.project) {
                    params.push(`project=${encodeURIComponent(newState.project)}`);
                }
                if (newState.dashboard !== undefined) {
                    params.push(`dashboard=${encodeURIComponent(newState.dashboard)}`);
                }
                if (newState.filter) {
                    params.push(
                        `filter=${encodeURIComponent(newState.filter[0])}:${encodeURIComponent(newState.filter[1])}`
                    );
                }
                const newSearch = `?${params.join('&')}`;
                if (newSearch !== window.location.search) {
                    window.history.pushState(newState, '', newSearch);
                }
                return newState;
            });
        },
        [setState]
    );

    return [state, setStateAndUrl];
};
