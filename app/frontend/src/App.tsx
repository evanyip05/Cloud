import React from "react";
import axios from 'axios';

export default function App() {
	return (
		<div>
	  		<button onClick={() => {
				axios.get("http://localhost:8080/test").then(res => {
					console.log(res)
				})
			}}>Mongo Put</button>
		</div>
 	);
}
