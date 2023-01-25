
import React, { useContext } from 'react';
import { Context } from './state'
import TreeNode from './TreeNode';

export default function Tree() {
    let { state } = useContext(Context);
    const entries = state['manifest']?.value?.filter((v: string) => v!=="")?.map((v: string) => "/" + v) ?? [];

    return (
        <div>
            <TreeNode entries={entries} path="" ></TreeNode>
        </div>
    )
};