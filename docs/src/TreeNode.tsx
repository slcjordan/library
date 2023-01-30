import React, { useState } from "react";
import { uniq } from "lodash";
import { NavLink } from "react-router-dom";

interface Props {
    entries: string[]
    path: string
};

export default function TreeNode(props: Props) {
    let name = props.path;
    let next = uniq(props.entries.map(curr => name + "/" + curr.substring(props.path.length + 1).split('/').shift()!));
    const [expanded, setExpanded] = useState(false);

    return (
        <div className="pl-2 ">
            <button className="text-left flex flex-row justify-start hover:text-sky-600" onClick={()=> setExpanded(!expanded)}>
                <span className="material-icons mr-2">{expanded ? "arrow_drop_down" : "arrow_right"}</span> {name === "" ? "root" : name}
            </button>
        {
            expanded
            ?
                next.map(curr => {
                let entries = props.entries.filter(name => name.startsWith(curr) && name !== curr);
                if (entries.length < 1) {
                    return (
                        <NavLink key={curr + "#link"} to={curr} className="hover:text-sky-600">
                            <div className="pl-2 text-left flex flex-row justify-start">
                            <span className="mr-2 material-icons-outlined">insert_drive_file</span>
                            {curr}
                            </div>
                        </NavLink>
                    )
                }
                return (<TreeNode
                        key={curr}
                        path={curr}
                        entries={entries}
                    ></TreeNode>)
                })
            :
                null
        }
        </div>
    )
};