# Objective:
Design an image hosting service where unregistered users can upload images.

## Functional requirements
* Unknown users can upload and reference images
* Should display number of views per image
* Resized versions of the images will be generated asynchronously 
and appear on the image page
* Images are not discoverable
* Images won't change
* No expiring
* Content moderation is not a feature
* Albums are not a feature
* Take-downs are not a feature

## Non-functional requirements
* Normal web size images (say, 6mb)
* We will not reject but rather compress images with size exceeding the limit and store compressed version as original
* 500k images uploads per day
* 5 million images read per day
* Growth 30-50% YoY
* 300ms p95 latency globally
* 99.0% availability for writes, 99.9% availability for reads
* Can use a cloud


## Estimations
### RPS:
```
       5          
 5 * 10           
 -------  ≈ 10 Upload RPS
  86400           
```
```
       6          
 5 * 10           
 -------  ≈ 100 Read RPS
  86400           
```


### Image storage:
<!-- http://www.sciweavers.org/free-online-latex-equation-editor -->
Let's assume that average "original" image size is 2MB.
In 1 year we will need `200 TB` just for storing original images:
```
        5                           8               
5  *  10   *  365  *  1MB ≈ 2  *  10  MB  =  200 TB 
```

Let's assume that we store 2 resized images: 10% and 50% of the original dimensions.  
Let's also assume that both resized images take `0.5MB` on average.  
That adds extra `50 TB` for resized versions, bringing total storage to `250 TB` in one year.

On S3 at the end of the first year that would cost approximately:
```
250_000GB * 0.022$/month ≈ 5000$ per month
```

Let's assume that we use our object storage with replication factor `3`.  
That means we need to store `750 TB` in one year.  
Assuming each server holds `10 TB`, we will need `75 servers`.  

### Bandwidth: 
Image uploads:
```
10 RPS * 2MB = 20 MB/sec = 160 Mbit/sec
```

Image downloads:
```
100 RPS * 2MB = 200 MB/sec = 1600 Mbit/sec
```

### Resizing images
Assuming that producing resized images takes 1 second,
we should ...

## Design
