"use client";
import Tree from "@/components/tree";
import { useEffect, useState } from "react";
import LiveTree from "@/components/live_tree";


const data = {
    tree: {
        name: "Dust",
        node_discovered: 2,
        children: [
            {
                name: "Earth",
                node_discovered: 0,
                children: []
            },
            {
                name: "Air",
                node_discovered: 1,
                children: []
            }
        ]
    },
    algorithm: "BFS",
    duration_ms: 0,
    visited_nodes: 1
};


export default function Home() {
  const [startAnimation, setStartAnimation] = useState(false);
  const [treeData, setTreeData] = useState(data.tree);
  const [inputValue, setInputValue] = useState('');
  const [resetAnimation, setResetAnimation] = useState(false);
  const [showLive, setShowLive] = useState(false);

  const [target, setTarget] = useState('');
  const [algorithm, setAlgorithm] = useState<'BFS' | 'DFS'>('BFS');
  const [mode, setMode] = useState<'single' | 'multiple'>('single');
  const [maxRecipes, setMaxRecipes] = useState<number>(30);

  const handleSubmit = async () => {
    const payload = {
      target,
      algorithm,
      mode,
      max_recipes: maxRecipes,
    };

    const res = await fetch(process.env.NEXT_PUBLIC_ENDPOINT as string, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    });

    const data = await res.json();
    setTreeData(data.tree);
    console.log('API Response:', data);
  };

  const handle_submit = () => {
    if (inputValue !== '') {
      try {
        const parsed = JSON.parse(inputValue);
        setTreeData(parsed);
        setInputValue('');
        console.log("changed the tree display!");
      } catch (err) {
        console.error('Invalid JSON:', err);
      }
    }
  }

  const onStartLive = () => {
    if (startAnimation) return;
    setStartAnimation(true);
    setShowLive(true);
  }

  const onStoppedLive = () => {
    console.log("Animation completed");
    setStartAnimation(false);
    setShowLive(false);
  }


  return (
    <div className="flex flex-col items-center w-full h-screen bg-gray-800 p-4 space-y-4">
      <h1 className="text-4xl">Recipe Tree</h1>
      <div className="relative w-4/5 h-[600px] overflow-auto border border-gray-600 rounded-lg">
        <div className="absolute min-w-full min-h-full flex justify-center items-start p-8">
          {
            !showLive ?
            <Tree node={treeData} type="root"/>
            :
            <LiveTree
            root={treeData}
            start={startAnimation}
            delay={800}
            onAnimationComplete={onStoppedLive}
            resetAnimation={resetAnimation}
            />
          }
        </div>
      </div>

        <div className="max-w-md mx-auto p-4 border rounded-xl shadow space-y-4 bg-white text-black">
        <input
          type="text"
          placeholder="Enter target (e.g., Dust)"
          value={target}
          onChange={(e) => setTarget(e.target.value)}
          className="w-full p-2 border rounded"
        />
          <div className="flex justify-between">
            <label>
              Algorithm:
              <select
                value={algorithm}
                onChange={(e) => setAlgorithm(e.target.value as 'BFS' | 'DFS')}
                className="ml-2 p-1 border rounded"
              >
                <option value="BFS">BFS</option>
                <option value="DFS">DFS</option>
              </select>
            </label>

            <label>
              Mode:
              <select
                value={mode}
                onChange={(e) => setMode(e.target.value as 'single' | 'multiple')}
                className="ml-2 p-1 border rounded"
              >
                <option value="single">Single</option>
                <option value="multiple">Multiple</option>
              </select>
            </label>
          </div>

          <input
            type="number"
            value={maxRecipes}
            onChange={(e) => setMaxRecipes(Number(e.target.value))}
            className="w-full p-2 border rounded"
            placeholder="Max Recipes"
            min={1}
          />

          <button
            onClick={handleSubmit}
            className="w-full bg-blue-600  py-2 rounded hover:bg-blue-700 text-white"
          >
            Search
          </button>
        </div>
        <button 
            onClick={onStartLive}
            className={`p-2 rounded  mt-4 ${startAnimation ? 'bg-gray-500' : 'bg-blue-500'}`}
            disabled={startAnimation}
            >
            Start Animation
        </button>
    </div>
  );
}