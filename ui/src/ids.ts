import React from 'react';

export const useIds = (project?: string): string[] => {
    const [ids, setIds] = React.useState<string[]>([]);

    React.useEffect(() => {
        if (project) {
            fetch(`./project/${project}/id`)
                .then((res) => res.json())
                .then(setIds);
        }
    }, [project]);

    return ids;
};
