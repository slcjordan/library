import React, { useEffect } from 'react';
import Dashboard from './Dashboard';
import Error from './Error';
import Markdown from './Markdown';
import Pill from './Pill';
import Search from './Search';
import Tree from './Tree';
import { Context, useDispatcher } from './state';
import { createBrowserRouter, RouterProvider } from "react-router-dom";


function App() {
  const dispatcher = useDispatcher();

  useEffect(() => {
      const key = "manifest";
      let state = dispatcher.state;
      if (!(state[key]?.loading) && state[key]?.value === undefined) {
          dispatcher.dispatch(key, fetch("/markdown/manifest.txt").then((response) => response.text()).then(body => body.split('\n')))
      }
  });
  const router = createBrowserRouter([
    {
      path: "/",
      element: <Dashboard></Dashboard>,
      errorElement:  <Error></Error>,
      children: [
        {
          path: "",
          element: <Markdown></Markdown>,
        },
        {
          path: "search",
          element: (<Search></Search>)
        },
        {
          path: "tree",
          element: (<Tree></Tree>)
        },
        {
          path: "cover.html",
          element: <iframe className="w-full h-full" title="coverage report" src="/test/cover.html"></iframe>,
        },
        {
          path: "*",
          element: <Markdown></Markdown>,
        },
      ],
    }
  ]);
  return (
    <Context.Provider value={dispatcher}>
      <RouterProvider router={router}></RouterProvider>
    </Context.Provider >
  )
}

export default App;