"use client";
import Tree from "@/components/tree";
import { useEffect, useState } from "react";
import LiveTree from "@/components/live_tree";


const treeData1 = {
  name: "Root",
  idx: 0,
  children: [
    {
      first: {
        name: "Child 1A",
        idx: 1,
      },
      second: {
        name: "Child 1B",
        idx: 2,
        children: [
          {
            first: {
              name: "Grandchild 1B-A",
              idx: 3,
            },
            second: {
              name: "Grandchild 1B-B",
              idx: 4,
            },
          }
        ]
      }
    },
    {
      first: {
        name: "Child 2A",
        idx: 5,
      },
      second: {
        name: "Child 2B",
        idx: 6,
      }
    }
  ]
};

const treeData2 = {
  name: "Root",
  idx: 0,
  children: [
    {
      first: {
        name: "Child 1A",
        idx: 1,
        children: [
          {
            first: {
              name: "Grandchild 1A-A",
              idx: 3,
              children: [
                {
                  first: { name: "Great-Grandchild 1A-A-1", idx: 7 },
                  second: { name: "Great-Grandchild 1A-A-2", idx: 8 }
                }
              ]
            },
            second: {
              name: "Grandchild 1A-B",
              idx: 4,
              children: [
                {
                  first: { name: "Great-Grandchild 1A-B-1", idx: 9 },
                  second: { name: "Great-Grandchild 1A-B-2", idx: 10 }
                }
              ]
            }
          }
        ]
      },
      second: {
        name: "Child 1B",
        idx: 2,
        children: [
          {
            first: { 
              name: "Grandchild 1B-A", 
              idx: 5,
              children: [
                {
                  first: { name: "Great-Grandchild 1B-A-1", idx: 11 },
                  second: { name: "Great-Grandchild 1B-A-2", idx: 12 }
                }
              ]
            },
            second: { 
              name: "Grandchild 1B-B", 
              idx: 6,
              children: [
                {
                  first: { name: "Great-Grandchild 1B-B-1", idx: 13 },
                  second: { name: "Great-Grandchild 1B-B-2", idx: 14 }
                }
              ]
            }
          }
        ]
      }
    },
    {
      first: {
        name: "Child 2A",
        idx: 15,
        children: [
          {
            first: { name: "Grandchild 2A-A", idx: 19 },
            second: { name: "Grandchild 2A-B", idx: 20 }
          }
        ]
      },
      second: {
        name: "Child 2B",
        idx: 16,
        children: [
          {
            first: { name: "Grandchild 2B-A", idx: 21 },
            second: { name: "Grandchild 2B-B", idx: 22 }
          }
        ]
      }
    },
    {
      first: {
        name: "Child 3A",
        idx: 17
      },
      second: {
        name: "Child 3B",
        idx: 18,
        children: [
          {
            first: { name: "Grandchild 3B-A", idx: 23 },
            second: { name: "Grandchild 3B-B", idx: 24 }
          }
        ]
      }
    }
  ]
};

export default function Home() {
  const [startAnimation, setStartAnimation] = useState(false);
  const [treeData, setTreeData] = useState(treeData1);
  const [inputValue, setInputValue] = useState('');
  const [resetAnimation, setResetAnimation] = useState(false);
  const [showLive, setShowLive] = useState(false);

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

  console.log(startAnimation);
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
      <div className="w-2/3 h-32 flex flex-row justify-center">
        <textarea
          value={inputValue}
          onChange={(e) => setInputValue(e.target.value)}
          placeholder="Input your JSON here!"
          className="bg-white w-9/10 h-32 rounded-2xl border-4 border-black text-black p-4 resize-none overflow-x-auto overflow-y-auto whitespace-pre text-wrap"
          spellCheck={false}
        />
        <button onClick={handle_submit}
          className={`w-1/10 h-full rounded-2xl border-black ml-4 text-xl ${startAnimation ? 'bg-gray-500' : 'bg-green-500'}`}
          disabled={startAnimation}>
          Enter
        </button>
      </div>
      <button 
          onClick={onStartLive}
          className={`p-2 rounded text-white mt-4 ${startAnimation ? 'bg-gray-500' : 'bg-blue-500'}`}
          disabled={startAnimation}
          >
          Start Animation
      </button>
    </div>
  );
}