let name = "Monkey";

let age = 1;

let inspirations = ["Scheme", "Lisp", "JavaScript", "Clojure"];

let book = [
	"title" : "Writing A Compiler in Go",
	"author" : "Thorsten Ball",
	"prequel" : "Writing An Interpreter in Go"
];

let printBookName = fn(book){
	let title = book["title"];
	let author = book["author"];
	puts(author + " - " + title);
};

printBookName(book);
// => prints: "Thorsten Ball - Writing A Compiler in Go"

let fibonacci = fn(x){
	if(x == 0){
		0
	}else{
		if(x == 1){
			return 1;
		}
		else{
			fibonacci(x-1) + fibnacci(x-2);
		}
	}
};

let map = fn(arr, f){
	let iter = fn(arr, accumulated){
		if(len(arr) == 0){
			accumulated
		}else{
			iter(rest(arr), push(accmulated, f(first(arr))));
		}
	};
	iter(arr, []);
};

let numbers = [1, 1+1, 4-1, 2*2, 2+3, 12/2];
map(numbers, fibonacci);
// => returns:[1,1,2,3,5,8]
