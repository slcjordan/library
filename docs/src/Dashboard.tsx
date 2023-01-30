import React from 'react';
import Pill from './Pill'
import { NavLink, Outlet } from "react-router-dom";
import LoadingOverlay from './LoadingOverlay';
import logo from './logo.png'

export default function Dashboard() {
    const navClass = ({ isActive }: { isActive: boolean }) => "rounded-br rounded-tl-lg " + (isActive ? "bg-white text-black hover:underline" : "hover:bg-white hover:text-black")
    return (
        <div className="flex flex-col md:flex-row font-semibold min-h-full min-w-full">
          <header className="flex flex-row md:flex-col pt-2 pl-4 bg-black text-white w-full md:w-2/12">
            <img src={logo} alt='' className="rounded-lg hidden md:flex self-center m-4 mt-8 w-1/2"></img>
            <NavLink className={navClass} to="/"><Pill>Home</Pill></NavLink>
            <NavLink className={navClass} to="/search"><Pill>Docs Search</Pill></NavLink>
            <NavLink className={navClass} to="/tree"><Pill>Docs Tree View</Pill></NavLink>
            <NavLink className={navClass} to="/cover.html"><Pill>Test Coverage</Pill></NavLink>
            <NavLink className={navClass} to="/api.md"><Pill>API Docs</Pill></NavLink>
          </header>
          <div className="border-t-2 md:border-t-0 md:border-l-2 border-black w-full md:w-10/12">
                <LoadingOverlay>
                    <Outlet></Outlet>
                </LoadingOverlay>
          </div>
        </div>
    )
}