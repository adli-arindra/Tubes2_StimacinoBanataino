"use client";
import Tree from "@/components/tree";

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
  return (
    <div className="flex flex-col items-center w-full h-screen bg-gray-800 p-4 space-y-4">
      <h1 className="text-4xl">Recipe Tree</h1>
      <div className="relative w-4/5 h-[500px] overflow-auto border border-gray-600 rounded-lg">
        <div className="absolute min-w-full min-h-full flex justify-center items-start p-8">
          <Tree node={treeData2} />
        </div>
      </div>
    </div>
  );
}