# Vibe Up Study

## Summary

Randomised controlled trials are the gold standard for the study of causal relationships, but they are often expensive to run, require large numbers of participants, and can take months or years to produce results. Integrating AI into the clinical trials landscape offers exciting potential to deliver faster results at lower costs based on smaller participant cohorts.

Vibe Up, which forms part of the three-year, multi-institutional Optimise Project, presents a novel opportunity to explore the utility of AI-driven response adaptive randomisation in a mental health context. This approach involves running a series of consecutive mini-trials; after each trial, the AI uses participant data to update its mathematical model based on its understanding of how well each intervention worked.

The Vibe Up trial is focused on addressing psychological distress among university students. Among this cohort, such distress can potentially lead to poor academic outcomes, discontinued study, and increased risk of suicide. Previous research shows that targeting mindfulness, physical activity and sleep through the use of psychological interventions can be effective in reducing distress; as such, these are the focus areas of the Vibe Up project.

## Useful Links

- [Black Dog Institute](https://www.blackdoginstitute.org.au/research-projects/vibe-up/)
- [Optimise](https://www.blackdoginstitute.org.au/research-centres/the-optimise-project/)

## Pulling Raw Data

> [!IMPORTANT]
> Users must be authenticated and authorised to access data. Speak to the administrator for more information.

Raw study data is stored as `csv` files on Google Cloud Storage (GCS). To obtain these files, use the [Google Cloud CLI](https://cloud.google.com/sdk/docs/install).

> [!NOTE]
> There are two versions of the CLI: `gsutil` and `gcloud`. The former is deprecated - use `gcloud` as it is faster and can handle multiple simultaneous downloads better.

```
gcloud storage cp "gs://<bucket-name>/<dir-name>/*.csv" .
```

I have purposely omitted the bucket name, etc. This command downloads the file to the current working directory (i.e,. `.`) but this path can be changed to point elsewhere.
Also note that:

1. The source path is enclosed in quotes as the glob can sometimes confuse the terminal.
2. There are options for recursive copying (`-r`), and continuing on error (`-c`).
3. `gcloud` automatically uses parallelism.
4. Gzipping can be done in flight, but only for uploads. This means that for downloads, compression must be done separately.
5. Globbing can be used to match `*` for files and `**` folders. Advanced pattern matching is not supported, however it can be achieved by first obtaining a list of files or folders in the parent directoy and then using pattern matching to filter the list. This list can then be used to pipe into the CLI for downloading.

Once the data has arrived, all the data is uncompressed. For large volume of data, [pigz](https://zlib.net/pigz/) is a good choice for fast compression.

```
# "pig-zee"
pigz -p 8 *.csv
```

