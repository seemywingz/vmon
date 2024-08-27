import React, { useState, useEffect } from "react";
import Pin from "./Pin";
import { Container, Fab } from "@mui/material";
import AddIcon from '@mui/icons-material/Add';

export default function Pins({ config }) {
    const [pins, setPins] = useState({}); // Initialize with an empty object

    useEffect(() => {
        if (config && config.pins) {
            setPins(config.pins);
        }
    }, [config]);

    const handleDeletePin = (pinNum) => {
        const updatedPins = { ...pins };
        delete updatedPins[pinNum];
        setPins(updatedPins);
    };

    const handleUpdatePin = (pinState, method) => {
        fetch(`http://${config.hostname}:${config.port}/api/pin`, {
            method: method,
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(pinState),
        }).then(response => {
            if (response.ok) {
                if (method === "DELETE") {
                    handleDeletePin(pinState.num);
                } else if (method === "POST") {
                    setPins(prevPins => ({
                        ...prevPins,
                        [pinState.num]: pinState,
                    }));
                }
            }
        }).catch(error => {
            console.error('Error updating pin:', error);
        });
    };

    if (!pins || Object.keys(pins).length === 0) {
        return null; // Return null if pins are not available or empty
    }

    // Map over the keys of the pins object
    const pinElements = Object.keys(pins).map((key) => {
        const pin = pins[key];
        return <Pin key={key} pinNum={key} props={pin} onUpdate={handleUpdatePin} />;
    });

    return (
        <Container sx={{
            display: 'flex',
            flexWrap: 'wrap',
            justifyContent: 'center',
        }}>
            <Fab size="small" color="primary" aria-label="add"
                sx={{
                    position: 'fixed',
                    top: '20px',
                    left: '20px',
                }}>
                <AddIcon />
            </Fab>
            {pinElements}
        </Container>
    );
}
