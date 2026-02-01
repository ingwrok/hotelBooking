import React from 'react';
import { cn } from '../../utils/cn.js';

const Button = ({
    children,
    variant = 'primary',
    size = 'md',
    className,
    fullWidth = false,
    ...props
}) => {

    const variants = {
        primary: "bg-primary text-white hover:bg-primary-light disabled:bg-primary/50",
        secondary: "bg-secondary text-white hover:bg-secondary/90",
        outline: "border-2 border-primary text-primary hover:bg-primary/5",
        ghost: "hover:bg-gray-100 text-gray-700",
        danger: "bg-red-500 text-white hover:bg-red-600"
    };

    const sizes = {
        sm: "px-3 py-1.5 text-sm",
        md: "px-5 py-2.5",
        lg: "px-8 py-3 text-lg"
    };

    return (
        <button
            className={cn(
                "rounded-lg font-medium transition-all duration-200 disabled:cursor-not-allowed flex items-center justify-center",
                variants[variant],
                sizes[size],
                fullWidth && "w-full",
                className
            )}
            {...props}
        >
            {children}
        </button>
    );
};

export default Button;
