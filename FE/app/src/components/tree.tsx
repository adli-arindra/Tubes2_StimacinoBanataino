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
};

export default function Tree({ node }: TreeProps) {
    return (
        <div className="flex flex-col-reverse items-center">
        <div className="flex flex-col items-center w-full">
            { node.children && node.children.length > 1 &&
                <>
                    <div className="w-full h-px bg-white"/>
                    <div className="w-px h-12    bg-white"/>
                </>
            }
            <div className="bg-green-500 text-white px-4 py-2 
                rounded-lg font-semibold shadow z-10 w-32 text-center">
                {node.name}
            </div>
        </div>

        {node.children && (
            <div className="flex space-x-8 relative pt-4">
                <div className="flex space-x-8">
                {node.children.map((child, index) => (
                    <div key={index} className="flex flex-col-reverse items-center">
                        <div className="w-px h-12 bg-white"/>
                        <div className="w-full h-px bg-white"/>
                        <div className="relative">
                            <div className="flex space-x-8">
                                <div className="flex flex-col items-center justify-end">
                                    <Tree node={child.first} />
                                    <div className="w-px h-24 bg-white"/>
                                    </div>
                                <div className="flex flex-col items-center justify-end">
                                    <Tree node={child.second} />
                                    <div className="w-px h-24 bg-white"/>
                                </div>
                            </div>
                        </div>
                    </div>
                ))}
                </div>
            </div>
        )}
        </div>
    );
}