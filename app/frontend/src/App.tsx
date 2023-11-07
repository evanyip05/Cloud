import React from "react";
import axios from 'axios';

export default function App() {
	return (
		<div style={{display:"flex", width:"100%", height:"100%", flexDirection:"row"}}>


	  		<div style={{display:"flex", width:"50%", height:"100%", flexDirection:"column", padding:"20px 10px 20px 20px"}}>
				<div>Input</div>
				<div style={{width:"100%", height:"auto", flexGrow:"1", paddingTop:"20px"}}>
				    <span style={{display:"block", width:"100%", height:"100%", padding:"5px", border:"1px solid black", borderRadius:"7.5px", overflowY:"auto"}} contentEditable></span>
				</div>
				<div style={{paddingTop:"20px"}}>
				    <button>add to db</button>
				    <button>find in db</button>
				</div>
			</div>


	  		<div style={{display:"flex", width:"50%", height:"100%", flexDirection:"column", padding:"20px 20px 20px 10px"}}>
				<div>Output</div>
				<div style={{width:"100%", height:"auto", flexGrow:"1", paddingTop:"20px"}}>
				    <span style={{display:"block", width:"100%", height:"100%", padding:"5px", border:"1px solid black", borderRadius:"7.5px", overflowY:"auto"}} contentEditable></span>
				</div>
			</div>
		</div>
 	);
}
