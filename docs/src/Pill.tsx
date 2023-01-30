import React from 'react';

interface Props {
    children?: React.ReactNode
}

export default function Pill(props: Props) {
    return (
        <div className='p-2'>
            {props.children}
        </div>
    )
}
