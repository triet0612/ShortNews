import {api_url} from "./index.js"
import { get } from "svelte/store"


export class Article {
  /**
   * @param {string} id 
   * @param {string} link 
   * @param {string} title 
   * @param {Date} pubdate 
   * @param {string} pubid 
   * @param {string} summary 
   */
  constructor(id, link, title, pubdate, pubid, summary) {
    this.id = id
    this.link = link
    this.title = title
    this.pubdate = pubdate
    this.pubid = pubid
    this.summary = summary
  }
}
/**
 * @param {number} page 
 * @returns {Promise<Article[]>}
 */
export async function articleFromApi(page) {
  /**
   * @type {Array<any>}
   */
  let jsResponse = await fetch(get(api_url) + `/articles?limit=5&start=${page*5}`)
    .then(res => res.json())
    .catch(err => {console.log(err);return []})
  /**
   * @type {Article[]}
   */
  let articles = []
  jsResponse.map(v => {
    articles.push(new Article(
      v["ArticleID"], v["Link"], v["Title"], v["PubDate"], v["PublisherID"], v["Summary"]
    ))
  })
  return articles
}
