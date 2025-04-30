"use client";
import { useEffect, useState } from "react";

type Recipe = {
    first: TreeNode;
    second: TreeNode;
};

type TreeNode = {
    name: string;
    idx: number;
    children?: Recipe[];
};

type TreeProps = {
    node: TreeNode;
    currentIdx: number;
    type?: string;
};

function Tree({ node, currentIdx, type }: TreeProps) {
    if (node.idx > currentIdx) return null;
    const isCurrent = node.idx === currentIdx;
    const nodeColor = isCurrent ? "bg-blue-500" : "bg-green-500";
    return (
        <div className="flex flex-col-reverse items-center">
            <div className="relative flex flex-col items-center w-full">
                {node.children && node.children.length > 1 && (
                    <>
                        <div className="w-full h-px bg-white" />
                        <div className="w-px h-12 bg-white" />
                    </>
                )}
                
                <div className={`${nodeColor} text-white px-4 py-2 
                    rounded-lg font-semibold shadow z-10 w-32 text-center`}>
                    {node.name}
                </div>
                {type === 'left' && (
                    <div className="absolute bottom-[-97px] left-1/2 w-1/2 h-1 bg-gray-800 -translate-x-full" />
                )}
                {type === 'right' && (
                    <div className="absolute bottom-[-97px] left-1/2 w-1/2 h-1 bg-gray-800 translate-x-0" />
                )}
            </div>
            {node.children && (
                <div className="flex space-x-8 relative pt-4">
                    <div className="flex space-x-8">
                        {node.children.map((child, index) => {
                            const showFirst = child.first.idx <= currentIdx;
                            const showSecond = child.second.idx <= currentIdx;

                            if (!showFirst && !showSecond) return null;

                            return (
                                <div key={index} className="flex flex-col-reverse items-center">
                                    {showFirst && showSecond && (
                                        <>
                                            <div className="w-px h-12 bg-white" />
                                            <div className="w-full h-px bg-white" />
                                        </>
                                    )}
                                    <div className="relative">
                                        <div className="flex space-x-8">
                                            <div className="flex flex-col items-center justify-end">
                                                <Tree node={child.first} currentIdx={currentIdx} type="left"/>
                                                {showFirst && <div className="w-px h-24 bg-white" />}
                                            </div>
                                            <div className="flex flex-col items-center justify-end">
                                                <Tree node={child.second} currentIdx={currentIdx} type="right"/>
                                                {showSecond && <div className="w-px h-24 bg-white" />}
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            );
                        })}
                    </div>
                </div>
            )}
        </div>
    );
}


type LiveTreeProps = {
    root: TreeNode;
    delay?: number;
    start: boolean;
    onAnimationComplete: () => void;
    resetAnimation: boolean;
};

export default function LiveTree({
    root,
    delay = 500,
    start,
    onAnimationComplete,
    resetAnimation,
}: LiveTreeProps) {
    const [currentIdx, setCurrentIdx] = useState(0);
    useEffect(() => {
        if (!start || resetAnimation) {
            setCurrentIdx(0);
        }

        if (!start || resetAnimation) return;

        const interval = setInterval(() => {
            setCurrentIdx((prevIdx) => {
                if (prevIdx >= getMaxIdx(root)) {
                    clearInterval(interval);
                    onAnimationComplete();
                    return prevIdx;
                }
                return prevIdx + 1;
            });
        }, delay);

        return () => clearInterval(interval);
    }, [start, root, delay, onAnimationComplete, resetAnimation]);

    return <Tree node={root} currentIdx={currentIdx} />;
}

function getMaxIdx(node: TreeNode): number {
    let maxIdx = node.idx;

    if (node.children) {
        for (const child of node.children) {
            maxIdx = Math.max(maxIdx, getMaxIdx(child.first), getMaxIdx(child.second));
        }
    }

    return maxIdx;
}