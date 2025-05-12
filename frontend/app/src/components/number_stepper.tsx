import React from "react";

interface NumberStepperProps {
    value: number;
    setValue: React.Dispatch<React.SetStateAction<number>>;
    min?: number;
    max?: number;
}

const NumberStepper = ({ value, setValue, min = 0, max = 10 } : NumberStepperProps) => {
    const decrement = () => {
        if (value > min) setValue(value - 1);
    };

    const increment = () => {
        if (value < max) setValue(value + 1);
    };

    return (
        <div style={{ display: "flex", alignItems: "center", gap: "0.5rem" }}>
        <button onClick={decrement} disabled={value <= min}>-</button>
        <span>{value}</span>
        <button onClick={increment} disabled={value >= max}>+</button>
        </div>
    );
};

export default NumberStepper;
