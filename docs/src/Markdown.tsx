
import React, { useContext, useEffect } from 'react';
import ReactMarkdown from 'react-markdown'
import { useLocation } from 'react-router-dom';
import { Context } from './state'
import remarkGfm from 'remark-gfm'
import './markdown.scss';

const keyPrefix = 'markdown-';

export default function Markdown() {
    let location = useLocation();
    let { state, dispatch } = useContext(Context);
    const key = keyPrefix + location.pathname;
    useEffect(() => {
        if (!(state[key]?.loading) && state[key]?.value === undefined) {
            dispatch(key, fetch("/markdown/" + (location.pathname === "/" ? "/README.md" : location.pathname)).then((response) => response.text()))
        }
    });
    return (
        <div className="markdown ml-4 p-2">
            <ReactMarkdown remarkPlugins={[remarkGfm]} >{ state[key]?.value ?? "" }</ReactMarkdown>
        </div>
    )
};