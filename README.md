# Efishery Test Assesement using Golang & Node JS by Riski Ramdan

# Clean Architecture principles

## Getting started (< 2mn)

```
git clone git@github.com:riskiramdan/efishery.git
cd efishery

docker-compose up //Create PostgresSQL & Redis Services
```
## Running Go Application 
```
cd golang                                 // Get in Golang Project root folder
cd databases                              // Get in Database Directory
rice embed-go                             // Generate seeder file
cd ..                                     // Back to main directory
go run cmd/efishery-migrate/main.go       // Run Database Migration
go run cmd/efishery-seeder/main.go        // Run Seeder 
go run cmd/efishery/main.go               // Run Go Application
```

## Running Node JS Application 
```
cd node                                   // Get in Node Project root folder
npm install                               // install all dependencies
npm start                                 // Run Node Application
```

## Postman Collection 
https://www.getpostman.com/collections/8b5a4199cefdf9b3a19a

## Domain Driven Architectures

Software design is a very hard thing. From years, a trend has appeared to put the business logic, a.k.a. the (Business) Domain, and with it the User, in the heart of the overall system. Based on this concept, different architectural patterns was imaginated. 

One of the first and main ones was introduced by E. Evans in its [Domain Driven Design approach](http://dddsample.sourceforge.net/architecture.html).

![DDD Architecture](/doc/DDD_architecture.jpg)

Based on it or in the same time, other applicative architectures appeared like [Onion Architecture](https://jeffreypalermo.com/2008/07/the-onion-architecture-part-1/) (by. J. Palermo), [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/) (by A. Cockburn) or [Clean Architecture](https://8thlight.com/blog/uncle-bob/2012/08/13/the-clean-architecture.html) (by. R. Martin).

This repository is an exploration of this type of architecture, mainly based on DDD and Clean Architecture, on a concrete and modern JavaScript application.
 
## DDD and Clean Architecture

The application follows the Uncle Bob "[Clean Architecture](https://8thlight.com/blog/uncle-bob/2012/08/13/the-clean-architecture.html)" principles and project structure :

### Clean Architecture layers

![Schema of flow of Clean Architecture](/doc/Uncle_Bob_Clean_Architecture.jpg)

### The Dependency Rule

> The overriding rule that makes this architecture work is The Dependency Rule. This rule says that source code dependencies can only point inwards. Nothing in an inner circle can know anything at all about something in an outer circle. In particular, the name of something declared in an outer circle must not be mentioned by the code in the an inner circle. That includes, functions, classes. variables, or any other named software entity.
  
src. https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html#the-dependency-rule
