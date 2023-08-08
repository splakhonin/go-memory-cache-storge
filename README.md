Cache Package:

The cache package provides a simple and efficient in-memory cache implementation for storing key-value pairs with optional expiration times. It includes a background cleaner to automatically remove expired entries from the cache.

Installation:

    go get - u github.com/splakhonin/go-memory-cache-storge

Usage
Creating a Cache
To create a new cache instance, you can use the NewCache function:

    c := cache.NewCache()
    defer c.StopCleaner() // Don't forget to stop the background cleaner when done
    c.Set("myKey", "myValue", 10*time.Second) // Entry will expire after 10 seconds

If no TTL is provided, the default TTL (5 seconds) will be used.

Retrieving Entries from the Cache
To retrieve an entry from the cache, use the Get method:

    value, err := c.Get("myKey")
    if err != nil {
    // Handle the error (key not found)
    } else {
    // Use the retrieved value

Deleting Entries from the Cache
You can delete entries from the cache using the Delete method:


    err := c.Delete("myKey")
    if err != nil {
    // Handle the error (key not found)
    } else {
    // Entry successfully deleted
    }

Cleaning Expired Entries
The cache automatically cleans up expired entries in the background. The cleanup interval is set to 5 seconds by default and can be adjusted by modifying the cleanInterval constant in the code. The cleaner can be stopped using the StopCleaner method:

    c.StopCleaner() // Stop the background cleaner

Manual Cache Cleaning
You can manually clean the cache using the CacheClean method to remove a specific entry by key:

    err := c.CacheClean("myKey")
    if err != nil {
    // Handle the error (key not found)
    } else {
    // Entry successfully removed
    }

Usage example
Here's a simple example of how to use the cache package:

    c := cache.NewCache()
    defer c.StopCleaner()
	c.Set("name", "John", 10*time.Second)
	c.Set("age", 30, 20*time.Second)

	name, _ := c.Get("name")
	age, _ := c.Get("age")

	fmt.Printf("Name: %s, Age: %d\n", name, age)

License
This cache package is open-source and available under the MIT License. Feel free to use, modify, and distribute it according to the terms of the license.