## Installation Notes

### Prerequisites
Configured MySQL database server on the machine


#### __Installation Steps__
* Download the latest release of AART API Server from [here](https://github.com/koreset/aart-api/releases). 

* Choose the executable/binary that's suitable for your operating system (Windows/Mac/Linux are currently supported) and save the file in an appropriate location on your server or machine.

* The API Server comes without a graphical user interface and interactions with the application can only happen via the command line window (Windows) or the Terminal (Mac/Linux)

* Open a console and navigate to the download location of the file

* Assuming that the filename of the download is api-0.9-beta.exe, run the following command at the prompt
```shell script
api-0.9-beta.exe --new-setup
```
This will start off a series of prompts for configuration information

```shell script
Running Initial Setup...
Enter your Database name:
```
Enter the name of the configured MySQL database you plan to use. This should ideally be an empty database. If the created database is ads, enter that name.

```shell script
Enter your Database Host:
``` 
Enter the database hostname. If the database is installed on the same machine as the server on which the application will be installed, enter "localhost"

```shell script
Enter your Database Port:
```
Provide the port on which the MySQL server listens on. The default is 3306

```shell script
Enter the Database Username:
```
The name of the database user that has the requisite read and write permissions to the database.

```shell script
Enter the Database Password:
```
The password for the database user.

```shell script
Enter the port for the Application (default 9090):
```
The API Server is RESTful application that needs to listen for requests on a port. The default is 9090. However, this can be configured to any port the application will need to listen on.

```shell script
Enter the IP or domain for the application (default localhost):
```

Supplying this information will result in a configuration for the application to run.

When the setup is complete, run the following command to start the server:
```shell script
api-0.9-beta.exe
```
This assumes the name of the file is api-0.9-beta.exe.
Ensure that the MySQL Server is up and running.

 
