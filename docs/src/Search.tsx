import levenshtein from 'fast-levenshtein';
import React, { useContext, useState } from 'react';
import { NavLink } from 'react-router-dom';
import { Context } from './state'

export default function Search() {
    let { state } = useContext(Context);
    const entries = state['manifest']?.value?.filter((v: string) => v!=="") ?? [];
    const [search, setSearch] = useState("");
    const [results, setResults] = useState<string[]>([]);

    // TODO add keypress navigation
    // TODO add focus handling
    return (
        <div key="one" className="p-4 flex-col flex items-start w-full mt-4">
            <form key="a" className="ml-8 w-8/12 border-2 border-black p-4 rounded-lg" >
                <div className="flex-row flex items-start justify-around">
                    <span className="material-icons text-4xl p-2">search</span>
                    <div className="w-11/12 mr-4 relative">
                    <div  className="absolute w-full flex flex-col">
                        <input key="b" type="search" className="p-2 border-2 border-black w-full"
                        onChange={
                            evt => {
                                evt.preventDefault();
                                let val = evt.target.value;
                                setSearch(val);
                                if (!val || val === "") {
                                    setResults([])
                                    return;
                                }
                                let curr: string[] = [...entries];
                                curr.sort((a: string, b: string) => levenshtein.get(a, val) - levenshtein.get(b, val))
                                curr.splice(10);
                                setResults(curr)
                            }
                        }
                        value={search}
                        ></input>
                        <ol className={results.length > 0 ? "bg-white border-black border-x-2 border-b-2" : ""}>
                            {results.map(curr => <li key={curr}><NavLink className="hover:bg-black hover:text-white p-2" to={"/" + curr}>{curr}</NavLink></li>)}
                        </ol>
                    </div>
                    </div>
                </div>
            </form>
        </div>
    )
};