import React from 'react';
import Pill from './Pill'
import { NavLink, } from "react-router-dom";
import noidea from './noidea.jpg'

export default function Error() {
      return  (
      <div className="flex flex-col items-center font-semibold min-h-full min-w-full justify-center p-8">
        <img src={noidea} alt="welp sorry" className="w-96"></img>
        <h1 className="text-center text-4xl">I have no idea what I am doing</h1>
        <NavLink to="/" className="rounded-br rounded-tl-lg bg-black text-white hover:underline m-2" ><Pill>Back Home</Pill></NavLink>
      </div>
      )
}