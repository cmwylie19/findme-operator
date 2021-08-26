import React, { useState } from "react";
import { GoogleMap, InfoWindow, Marker } from "@react-google-maps/api";

const divStyle = {
  background: `white`,
};

function Map({ markers, center }) {
  const [activeMarker, setActiveMarker] = useState(null);

  const handleActiveMarker = (marker) => {
    if (marker === activeMarker) {
      return;
    }
    setActiveMarker(marker);
  };

  const handleOnLoad = (map) => {
    const bounds = new window.google.maps.LatLngBounds();
    markers.forEach(({ position }) => bounds.extend(position));
    map.fitBounds(bounds);
  };

  return (
    <>
      {center[0] !== "--" && (
        <GoogleMap
          //   onLoad={handleOnLoad}
          zoom={18}
          center={{
            lat: center[0],
            lng: center[1],
          }}
          onClick={() => setActiveMarker(null)}
          mapContainerStyle={{ width: "100%", height: "100vh" }}
        >
          {markers.map(({ name, position }, i) => (
            <Marker
              key={i}
              position={position}
              onClick={() => handleActiveMarker(i)}
            >
              {activeMarker === i ? (
                <InfoWindow
                  position={position}
                  onCloseClick={() => setActiveMarker(null)}
                >
                  <div style={divStyle}>{name}</div>
                </InfoWindow>
              ) : null}
            </Marker>
          ))}
        </GoogleMap>
      )}
    </>
  );
}

export default Map;
