import React, { useContext } from 'react';
import { Context } from './state'
import { find } from 'lodash'

interface Props {
    children?: React.ReactNode
}

export default function LoadingOverlay(props: Props) {
    let { state } = useContext(Context);
    const loading = !!find(Object.values(state??{}), ['loading', true])
    // TODO add spinner
    return (
        <div className={("h-full w-full" + (loading ? "bg-gray-50 opacity-80" : ""))}>
            {props.children}
        </div>
    )
}
