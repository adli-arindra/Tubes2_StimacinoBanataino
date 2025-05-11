"use client";
import { useEffect, useState } from "react";

type TreeNode = {
    name: string;
    node_discovered: number;
    children?: TreeNode[];
};

type TreeProps = {
    node: TreeNode;
    currentIdx: number;
    type?: string;
};

function Tree({ node, currentIdx, type }: TreeProps) {
    if (node.node_discovered > currentIdx) return null;
    const isCurrent = node.node_discovered === currentIdx;
    const nodeColor = isCurrent ? "bg-blue-500" : "bg-green-500";

    const showFirst = node.children && node.children.length > 0 
        ? node.children[0].node_discovered <= currentIdx 
        : false;

    const showSecond = node.children && node.children.length > 1 
        ? node.children[1].node_discovered <= currentIdx 
        : false;

    return (
        <div className="flex flex-col-reverse items-center">
            <div className="relative flex flex-col items-center w-full">
                {node.children && node.children.length > 1 && (
                    <>
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
                        {(showFirst || showSecond) && (
                            <div className="flex flex-col-reverse items-center">
                                {showFirst && showSecond && (
                                    <>
                                        <div className="w-px h-12 bg-white" />
                                        <div className="w-full h-px bg-white" />
                                    </>
                                )}
                                <div className="relative">
                                    <div className="flex space-x-8">
                                        <div className="flex flex-col items-center justify-end">
                                            <Tree node={node.children[0]} currentIdx={currentIdx} type="left" />
                                            {showFirst && <div className="w-px h-24 bg-white" />}
                                        </div>
                                        <div className="flex flex-col items-center justify-end">
                                            <Tree node={node.children[1]} currentIdx={currentIdx} type="right" />
                                            {showSecond && <div className="w-px h-24 bg-white" />}
                                        </div>
                                    </div>
                                </div>
                            </div>
                        )}
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
    let maxIdx = node.node_discovered;

    if (node.children) {
        for (const child of node.children) {
            maxIdx = Math.max(maxIdx, getMaxIdx(node.children[0]), getMaxIdx(node.children[1]));
        }
    }

    return maxIdx;
}