import React, { useEffect, useState } from "react";
import { useLoadScript } from "@react-google-maps/api";
import Map from "./Map";

const divStyle = {
  display: `flex`,
  justifyContent: `center`,
  alignItems: `center`,
  height: `100vh`,
  width: `100%`,
  padding: 1,
};

function App() {
  const [coords, setCoords] = useState([0.0, 0.0]);

  const { isLoaded } = useLoadScript({
    googleMapsApiKey: "", // Add your API key
  });

  var options = {
    enableHighAccuracy: true,
    timeout: 5000,
    maximumAge: 0,
  };

  function success(pos) {
    const { latitude, longitude, accuracy } = pos.coords;
    setCoords([parseFloat(latitude), parseFloat(longitude)]);

    console.log("Your current position is:");
    console.log(`Latitude : ${latitude}`);
    console.log(`Longitude: ${longitude}`);
    console.log(`More or less ${accuracy} meters.`);
  }

  function error(err) {
    console.warn(`ERROR(${err.code}): ${err.message}`);
  }

  useEffect(() => {
    navigator.geolocation.getCurrentPosition(success, error, options);
  }, []);
  return coords[0] === 0.0 ? (
    <div style={divStyle}>
      You need to allow access to Geolocation to use this app.
    </div>
  ) : (
    <>
      {isLoaded ? (
        <Map
          markers={[
            {
              name: `Me! ${coords[0]} ${coords[1]} `,
              position: {
                lat: coords[0],
                lng: coords[1],
              },
            },
          ]}
          center={[parseFloat(coords[0]), parseFloat(coords[1])]}
        />
      ) : null}
    </>
  );
}

export default App;
