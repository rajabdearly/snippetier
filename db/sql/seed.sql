INSERT INTO users (username, email, full_name)
VALUES ('some_username', 'user@example.com', 'Semashka Vicini');


-- Inserting data into the "snippets" table with random code snippets
INSERT INTO snippets (name, description, content, user_id) VALUES
    ('JavaScript Greeting', 'A JavaScript function to greet a person', 'function greet(name) {
    console.log("Hello, '' + name + ''!");
    }

    greet("John");', 1),

    ('Looping in JavaScript', 'A JavaScript loop example', 'for (let i = 1; i <= 5; i++) {
    console.log("Iteration " + i);
    }',1),

    ('Java Hello World', 'A Java program to print "Hello, World!"', 'public class HelloWorld {
    public static void main(String[] args) {
        System.out.println("Hello, World!");
    }
    }', 1),

    ('Calculating Sum in JavaScript', 'JavaScript code to calculate the sum of an array', 'const numbers = [1, 2, 3, 4, 5];
    const sum = numbers.reduce((acc, val) => acc + val, 0);
    console.log("Sum:", sum);', 1),

    ('Python Factorial', 'Python code to calculate the factorial of a number', 'def factorial(n):
    if n == 0:
        return 1
    else:
        return n * factorial(n-1)

    result = factorial(5)
    print("Factorial of 5 is:", result);', 1);
